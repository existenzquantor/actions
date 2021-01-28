
effect('Pull', [not(left)], [left]).
effect('Pull', [left], [not(left)]).
effect('Wait', [], []).
effect(0, [not(left)], [dead5]).
effect(0, [left], [dead1]).

init([not(dead5), not(dead1), not(left)]).
plan('Pull').
goal([not(dead5)]).