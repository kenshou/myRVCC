#!/bin/bash
assert () {
  # 程序运行的 期待值 为参数1
    expected="$1"
    # 输入值 为参数2
    input="$2"
    ./build/rvcc "$input" > build/tmp.s || exit
    riscv64-unknown-elf-gcc -static build/tmp.s -o build/tmp
    spike pk build/tmp

    actual="$?"
    if [ "$actual" = "$expected" ]; then
      echo "$input => $actual: Pass"
    else
      echo "$input => $expected expected, but got $actual"
      exit 1
    fi
}

# assert 期待值 输入值
# [1] 返回指定数值
assert 0 0
assert 42 42

# [2] 支持+ -运算符
assert 34 '12-34+56'

# [3] 支持空格
assert 41 ' 12 + 34 - 5 '
echo OK