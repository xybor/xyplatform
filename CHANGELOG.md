# Developing

1.  Add FileEmitter and RotatingFileEmitter to xylog.
2.  Xylog now can be easy to use key-value fields and extra values.
3.  Add Codacy analysis.

# V1.0.1

1.  Fix bugs.
2.  Refactor the design of xylog Handler.
3.  Add github workflow to test.
4.  Write unittest for all modules.
5.  Xycond asserts a Error instead of string.

# V1.0.0

This release completed the following libraries:

1.  xycond supports to check many types of condition and panic if the condition
fails.

2.  xyerror contains special errors that are good for error comparison and
debugging.

3.  xylock contains wrapper structs of built-in sync library, such as
`sync.Mutex` or `semaphore.Weighted`.

4.  xylog provides flexible logging methods to the program.

5.  xysched provides a mechanism of job scheduling in future with a simple
    syntax.

6.  xyselect is a library used to call `select` with an unknown number of `case`
    statements.
