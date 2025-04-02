"""
Module of basic tools for generating testing simple tasks.
"""

from cProfile import Profile
from copy import deepcopy
from dataclasses import dataclass
from itertools import combinations, permutations
from pstats import SortKey, Stats
from random import randint, seed, shuffle, choice
from abc import ABCMeta, abstractmethod
import time
from typing import Any, Callable, Generator, List

from tqdm import tqdm


class BColors:
    """
    Colors for highlighting special words.
    """

    HEADER = "\033[95m"
    OKBLUE = "\033[94m"
    OKCYAN = "\033[96m"
    OKGREEN = "\033[92m"
    WARNING = "\033[93m"
    FAIL = "\033[91m"
    ENDC = "\033[0m"
    BOLD = "\033[1m"
    UNDERLINE = "\033[4m"


@dataclass
class TestCase:
    """
    Test case.
    """

    args: List[Any]
    possible_resulsts: List[Any]

    def check(self, answer: Any) -> bool:
        """
        Whether answer in possible results.
        """

        return answer in self.possible_resulsts


def test(func, test_cases: List[TestCase]):
    """
    Simple test func.
    """

    for iters, test_case in enumerate(test_cases):
        res = func(test_case)

        if not test_case.check(res):
            print(
                f"Wrong answers in test {iters+1}, "
                f"got {res}, "
                f"expected one of {test_case.possible_resulsts}"
            )


def all_sub_lists(lst: list):
    """
    Return all possible divides of target list on
    lists. For example [1, 2, 3] -> [
        [[1], [2], [3]],
        [[1, 2], [3]],
        [[1], [2, 3]],
        [[1, 2, 3]]
    ].
    """

    if len(lst) < 2:
        return [lst]

    for perm in permutations(lst):
        yield [perm]
        for borders_count in range(1, len(lst)):
            for borders in combinations(range(1, len(lst)), r=borders_count):
                new_lst = [perm[: borders[0]]]
                last_b = borders[0]

                for b in borders[1:]:
                    new_lst.append(perm[last_b:b])
                    last_b = b

                low_border = borders[-1]
                new_lst.append(perm[low_border:])

                yield new_lst


class GT(metaclass=ABCMeta):
    """
    Base class for any generating argument. The basic idea is
    that arguments generate their values on demand, but
    return the same value within a single test.
    """

    @abstractmethod
    def generate(self):
        """
        Generating some value.
        """

        raise NotImplementedError

    @property
    @abstractmethod
    def val(self):
        """
        Returns lastly generated value.
        """

        raise NotImplementedError


class GChar(GT):
    """
    Generates one char of available.
    """

    def __init__(self, available_letters: str) -> None:
        self.al = available_letters
        self.__val = None

    @property
    def val(self):
        if self.__val is None:
            self.__val = self.generate()

        return self.__val

    @val.setter
    def val(self, value):
        self.__val = value

    def generate(self):
        self.val = choice(self.al)

        return self.val


class GFrozStr(GT):
    """
    Generates shuffle of str.
    """

    def __init__(self, cur_ctr: str) -> None:
        self.al = cur_ctr
        self.__val = None

    @property
    def val(self):
        if self.__val is None:
            self.__val = self.generate()

        return self.__val

    @val.setter
    def val(self, value):
        self.__val = value

    def generate(self):
        cur = list(self.al)
        shuffle(cur)
        self.val = "".join(cur)

        return self.val


class GInt(GT):
    """
    Generates int from / up to int or another GInt.
    """

    def __init__(self, min_v: int | GT, max_v: int | GT) -> None:
        self.min = min_v
        self.max = max_v
        self.__val = None

    @property
    def val(self):
        if self.__val is None:
            self.__val = self.generate()

        return self.__val

    @val.setter
    def val(self, value):
        self.__val = value

    def generate(self):
        val_min = self.min.val if isinstance(self.min, GT) else self.min
        val_max = self.max.val if isinstance(self.max, GT) else self.max

        self.val = randint(val_min, val_max)

        return self.val


class GUInt(GInt):
    """
    GInt with min value 1.
    """

    def __init__(self, max_v: int | GT) -> None:
        super().__init__(1, max_v)


class GList(GT):
    """
    Generates list of GT.
    """

    def __init__(
        self,
        cur_type: GT,
        amount: int | GUInt,
        *args,
        list_func: Callable = list,
    ) -> None:
        self.type = cur_type
        self.amount = amount
        self.args = args
        self.list_func = list_func
        self.__val = None

    @property
    def val(self):
        if self.__val is None:
            self.__val = self.generate()
        return self.__val

    @val.setter
    def val(self, value):
        self.__val = value

    def generate(self):
        amount = (
            self.amount.val if isinstance(self.amount, GT) else self.amount
        )
        cur = [self.type.generate(*self.args) for _ in range(amount)]

        shuffle(cur)
        self.val = self.list_func(cur)

        return self.val


class GTuple(GT):
    """
    Generates Tuple of GT.
    """

    def __init__(self, gargs: List[GT]) -> None:
        self.gargs = gargs
        self.__val = None

    @property
    def val(self):
        if self.__val is None:
            self.__val = self.generate()

        return self.__val

    @val.setter
    def val(self, value):
        self.__val = value

    def generate(self):
        self.val = tuple(garg.val for garg in self.gargs)

        return self.val


class GTester:
    """
    Operating tests fo GT. all_args - list of args, that
    need to be generated in a sequence, than return args only
    of part of them. universal_function - one that
    is probably unefficient, but will return right answer
    for testing results from func.
    """

    def __init__(
        self,
        func: Callable,
        universal_func: Callable,
        func_args: List[GT],
        all_args: List[GT],
    ) -> None:
        self.func = func
        self.ufunc = universal_func
        self.func_args = func_args
        self.all_args = all_args

    def generate_1(self) -> TestCase:
        """
        Generates values in all args, then return only func args.
        """

        for arg in self.all_args:
            arg.generate()

        return TestCase(
            args=[arg.val for arg in self.func_args], possible_resulsts=[]
        )

    def generate_n(self, amount: int) -> Generator[TestCase, Any, None]:
        """
        Calles generate_1 from 1 to infinite times.
        If amount less than 0 - generates infinitly.
        """

        i = 0
        while i < amount:
            seed(i)
            yield self.generate_1()
            i += 1

    def test(
        self,
        amount: int,
        time_limit: float = 1.0,
        print_right: int = 0,
        fail_on: int = None,
        dont_count_answers: list = None,
    ):
        """
        This test checking working of function return
        with small arguments. But it should not test
        corner case when a = 10 ** 9 or something.

        If amount is -1 then it creates endless tests.
        With flag fail_on - fails on n failed test.
        """

        if not dont_count_answers:
            dont_count_answers = []

        failed = 0
        arg_gen = self.generate_n(amount)

        with tqdm() as pbar:
            for iters, test_case in enumerate(arg_gen):
                failed += 1
                prev_test_case = deepcopy(test_case)

                start = time.perf_counter()
                res = self.func(test_case)
                end = time.perf_counter()

                res_time = end - start
                str_time = f"{BColors.BOLD}TIME: {res_time} s{BColors.ENDC}"
                test_case.possible_resulsts = [self.ufunc(prev_test_case)]

                if res not in test_case.possible_resulsts:
                    print(
                        f"\n{BColors.FAIL}=============TEST "
                        f"{BColors.OKCYAN}#{iters}{BColors.FAIL}"
                        f" FAILED=============\nArguments: {test_case}"
                        f"\n{str_time}{BColors.ENDC}\n"
                        f"EXPECT:\t{test_case.possible_resulsts}\nGOT:\t{res}"
                        f"{BColors.ENDC}"
                    )
                elif res_time > time_limit:
                    print(
                        f"\n{BColors.WARNING}=============TEST "
                        f"{BColors.OKCYAN}#{iters}{BColors.WARNING}"
                        f" TIME LIMIT EXCEEDED=============\n"
                        f"{BColors.OKGREEN}Answer is right\n{BColors.WARNING}"
                        f"Arguments: {test_case}"
                        f"\n{str_time}{BColors.ENDC}\n"
                        f"EXPECT:\t{test_case.possible_resulsts}\nGOT:\t{res}"
                        f"{BColors.ENDC}"
                    )
                elif res in test_case.possible_resulsts:
                    failed -= 1

                    if res in dont_count_answers:
                        continue

                    if print_right:
                        print(
                            f"\n{BColors.OKGREEN}=============TEST "
                            f"{BColors.OKBLUE}#{iters}{BColors.OKGREEN}"
                            f" PASSED=============\n{str_time}{BColors.ENDC}"
                        )

                        if print_right == 2:
                            print(f"Arguments: {test_case}\nGOT:\t{res}")

                if fail_on is not None and failed >= fail_on:
                    print(
                        f"\n{BColors.FAIL}=============TOTAL {failed}"
                        f" FAILS EXCEEDED============={BColors.ENDC}"
                    )

                    return

                pbar.update(1)

    def test_profile(self, amount: int, time_limit: float = 1.0):
        """
        Testing function on TL error.
        """

        arg_gen = self.generate_n(amount)

        with tqdm() as pbar:
            for iters, test_case in enumerate(arg_gen):
                with Profile() as profile:
                    self.func(test_case)
                    profile_result = Stats(profile)

                if profile_result.total_tt > time_limit:
                    print(
                        f"\n{BColors.WARNING}=============TEST "
                        f"{BColors.OKCYAN}#{iters}{BColors.WARNING}"
                        f" TIME LIMIT EXCEEDED============={BColors.ENDC}\n"
                        f"Total {BColors.WARNING}{profile_result.total_tt}"
                        f"{BColors.ENDC} more than "
                        f"{BColors.OKGREEN}{time_limit}{BColors.ENDC}"
                        f"{BColors.ENDC}"
                    )
                    profile_result.strip_dirs().sort_stats(
                        SortKey.TIME, SortKey.CALLS
                    ).print_stats()

                pbar.update(1)


def simple_example():
    """
    Simple example for module entities.
    """

    n = GUInt(10)
    m = GInt(-10, 10)
    B = GTuple([n, m])  # pylint: disable=C0103:invalid-name
    A = GList(m, n)  # pylint: disable=C0103:invalid-name

    def f(test_case: TestCase = None):
        if test_case:
            # pylint: disable=W0612:unused-variable
            n, m = test_case.args[0]
            # pylint: disable=C0103:invalid-name
            X = test_case.args[1]
        else:
            # pylint: disable=W0612:unused-variable
            n, m = [int(x) for x in input().split()]
            # pylint: disable=C0103:invalid-name
            X = [int(x) for x in input().split()]

        res = 0

        for xx in X:
            # little logic error
            if xx != 2:
                res += xx
            if xx == -10:
                time.sleep(1.3)

        return res

    def uf(test_case: TestCase = None):
        if test_case:
            # pylint: disable=W0612:unused-variable
            n, m = test_case.args[0]
            # pylint: disable=C0103:invalid-name
            X = test_case.args[1]
        else:
            # pylint: disable=W0612:unused-variable
            n, m = [int(x) for x in input().split()]
            # pylint: disable=C0103:invalid-name
            X = [int(x) for x in input().split()]

        return sum(X)

    tester = GTester(f, uf, [B, A], [n, m, B, A])
    tester.test(10, print_right=1)


if __name__ == "__main__":
    simple_example()
