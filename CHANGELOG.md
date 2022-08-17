# Developmet
1. Fix ClosedChannelError in xyselect.

# V1.0.0
This release completed the following libraries:
1. xycond supports to check many types of condition and panic if the condition
fails.
2. xyerror contains special errors that are good for error comparison and
debugging.
3. xylock contains wrapper structs of built-in sync library, such as
`sync.Mutex` or `semaphore.Weighted`.
4. xylog provides flexible logging methods to the program.
5. xysched provides a mechanism of job scheduling in future with a simple
syntax.
6. xyselect is a library used to call `select` with an unknown number of `case`
statements.
