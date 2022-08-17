# xyplatform
[![Xybor founder](https://img.shields.io/badge/xybor-huykingsofm-red)](https://github.com/huykingsofm)
![GitHub Repo stars](https://img.shields.io/github/stars/xybor/xyplatform?color=yellow)
![GitHub top language](https://img.shields.io/github/languages/top/xybor/xyplatform?color=lightblue)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/xybor/xyplatform)

Xyplatform contains platform libraries developed by Xybor.

# List of libraries
1. [xycond](./xycond) supports to check many types of condition and panic if the
condition fails.
2. [xyerror](./xyerror) contains special errors that are good for error
comparison and debugging.
3. [xylock](./xylock) contains wrapper structs of built-in `sync` library, such
as `sync.Mutex` or `semaphore.Weighted`.
4. [xylog](./xylog) provides flexible logging methods to the program.
5. [xyselect](./xyselect) is a library used to call `select` with an unknown
number of `case` statements.
