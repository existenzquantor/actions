:- module(action_helpers, [without_last/2, without_first_two/2, bash_command/2, empty_once/1]).

without_last([_], []).
without_last([X|Xs], [X|WithoutLast]) :- 
    without_last(Xs, WithoutLast).

without_first_two([_, _ | L], L).

bash_command(Command, Output) :-
    process_create(path(bash),
            ['-c', Command],
            [stdout(pipe(Out))]),
    read_string(Out, _, Output),
    close(Out).

empty_once(P) :-
    findall(X, (member(X, P), X = empty), L),
    length(L, 1).