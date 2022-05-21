#!/usr/bin/python3

# find all non negative integer solutions for a^2 + b^2 = 5^2.
# cf. https://stackoverflow.com/a/13398853

import z3

a = z3.Int('a')
b = z3.Int('b')

s = z3.Solver()
s.add(z3.And(a >= 0, b >= 0))
s.add(a * a + b * b == 5 * 5)

while s.check() == z3.sat:
    m = s.model()
    print(m)
    s.add(z3.Or(a != m[a], b != m[b]))
