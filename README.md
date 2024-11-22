# myRVCC
参考[rvcc](https://github.com/sunshaoce/rvcc)使用go来实现一个riscv的c编译器

# macos上的环境配置
添加新源
brew tap riscv-software-src/riscv  
安装 riscv-tools工具链  
brew install riscv-tools  
brew install riscv-pk  
brew install riscv-isa-sim

riscv64-unknown-elf-gcc -static tmp.s -o tmp  
spike pk tmp
