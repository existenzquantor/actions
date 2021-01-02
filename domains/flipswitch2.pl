effect('FlipSwitch', [on], [not(on)]).
effect('FlipSwitch', [not(on)], [on, yay]).

init([not(on), not(yay)]).
goal([not(on)]).
plan('FlipSwitch':'FlipSwitch').

