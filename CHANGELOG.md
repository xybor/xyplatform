# V0.0.3 (Aug 30, 2022)

1.  Add FileEmitter and RotatingFileEmitter to xylog.
2.  Add EventLogger to xylog.
3.  Add report analysis badges.
4.  Reduce logging time (3x faster).
5.  Refactor xycond to be shorter and more readable.
6.  Scheduler now can be identified by name.
7.  Add log to xysched and xyselect.
8.  Future now can early stop.
9.  Methods with formatting string are splitted to two seperated methods (with f
    and non-f suffix).

# V0.0.2 (Aug 21, 2022)

1.  Fix bugs.
2.  Refactor the design of xylog Handler.
3.  Add github workflow to test.
4.  Write unittest for all modules.
5.  Xycond asserts a Error instead of string.

# V0.0.1 (Aug 17, 2022)

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
