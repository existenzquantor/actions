:- module(actions_core, [classify_actions/1, classify_plan/1, prepare_owl/1, names/2, contexts/4, causedFacts/3, causedFacts/1, reasons/3]).
:- use_module("owl.pl", [prepare_owl/1]).
:- use_module("helpers.pl", [bash_command/2, without_last/2, without_first_two/2, empty_once/1]).
:- use_module("../../causality/core/interpreter.pl", [do/3, action/1, finally/2]).
:- use_module("../../causality/core/programs.pl", [program_to_list/2]).
:- use_module("../../causality/core/causality.pl", [cause_empty_temporal/3, reason_empty_temporal/4,reason_empty_temporal_nogoal/4]).

classify_actions(L) :-
    plan(Plan0),
    program_to_list(Plan0, Plan),
    classify_actions(Plan, [], L).
classify_actions([], L, L).
classify_actions([_ | R], L, E) :-
    length(L, N),
    classify_action(N, AL),
    append(L, [AL], L2),
    classify_actions(R, L2, E).

classify_action(N, L) :-
    format(atom(S), "java -jar ./reasoner/HermiT.jar -S:Action~w ./temp/temp.owl", [N]),
    bash_command(S, O),
    extract_answer_from_hermit(O, L).

classify_plan(L) :-
    S = "java -jar ./reasoner/HermiT.jar -S:Plan0 ./temp/temp.owl",
    bash_command(S, O),
    extract_answer_from_hermit(O, L).

names(P, L) :-
    program_to_list(P, L).

contexts(A, C, Cs, E) :-
    action(A),
    append(Cs, [C], E).
contexts(A : R, C, Cs, E) :-
    do(A, C, Cnext),
    append(Cs, [C], Cs2),
    contexts(R, Cnext, Cs2, E).

causedFacts(N, Program, Facts):-
    program_to_list(Program, PL),
    nth0(N, PL, Action),
    findall(Reason, (reason_empty_temporal_nogoal(Reason, Action, Program, Witness), program_to_list(Witness, WL), nth0(N, WL, empty), empty_once(WL)), F),
    sort(F, Facts).
causedFacts(Facts) :-
    plan(P),
    findall(X, (finally(P, X),cause_empty_temporal(P, X, _)), F),
    sort(F, Facts).
    

reasons(N, Program, Facts):-
    program_to_list(Program, PL),
    nth0(N, PL, Action),
    findall(Reason, (reason_empty_temporal(Reason, Action, Program, Witness), program_to_list(Witness, WL), nth0(N, WL, empty), empty_once(WL)), F),
    sort(F, Facts).

extract_answer_from_hermit(O, L) :-
    split_string(O, "\n", "", L0),
    without_last(L0, L1),
    without_first_two(L1, L2),
    extract_actions(L2, [], L).

extract_actions([], L, L).
extract_actions([A | L], B, E) :-
    split_string(A, ":", "", AL),
    nth0(1, AL, AA),
    append(B, [AA], L2),
    extract_actions(L, L2, E).