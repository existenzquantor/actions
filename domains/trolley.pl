
effect(pull, [not(left)], [left]).
effect(pull, [left], [not(left)]).
effect(wait, [], []).
effect(0, [not(left)], [dead5]).
effect(0, [left], [dead1]).

init([not(dead5), not(dead1), not(left)]).
plan(pull).
goal([not(dead5)]).