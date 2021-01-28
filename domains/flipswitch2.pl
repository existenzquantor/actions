effect('FlipSwitch', [on], [not(on)]).
effect('FlipSwitch', [not(on)], [on]).

init([not(on)]).
goal([not(on)]).
plan('FlipSwitch':'FlipSwitch').

