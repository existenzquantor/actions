:- module(actions_owl, [prepare_owl/1]).
:- use_module("core.pl", [names/2, contexts/4, causedFacts/3, reasons/3]).
:- use_module("helpers.pl", [without_last/2]).

prepare_owl(OntologyPath) :-
    plan(Plan),
    to_owl(Plan, S),
    append_owl(OntologyPath, S, New),
    write_owl(New).

lits_to_owl_strings(not(A), S, F) :-
    format(atom(X), " ObjectComplementOf(:~w)", A),
    string_concat(S, X, F),!.
lits_to_owl_strings(A, S, F) :-
    \+is_list(A),
    format(atom(X), " :~w", A),
    string_concat(S, X, F),!.
lits_to_owl_strings([], S, SF) :-
    format(atom(SF), " ObjectIntersectionOf(~w owl:Thing)", S).
lits_to_owl_strings([A|R], S, F) :-
    lits_to_owl_strings(A, S, F2),
    lits_to_owl_strings(R, F2, F).

to_owl(Plan, S) :-
    names(Plan, Names),
    length(Names, N),
    to_owl(Plan, 0, N, "", S).
to_owl(_, N, N, S, S).
to_owl(Plan, N, Nmax, Str, S) :-
    N < Nmax,
    to_owl(Plan, N, O),
    string_concat(Str, O, Str2),
    N2 is N + 1,
    to_owl(Plan, N2, Nmax, Str2, S).
to_owl(Plan, N, S) :-
     % Names
    names(Plan, Names),
    nth0(N, Names, Action),
    lits_to_owl_strings(Action, "", SName),
    % Context
    init(C),
    contexts(Plan, C, [], Contexts),
    nth0(N, Contexts, Context),
    lits_to_owl_strings(Context, "", SContexts),
    % Caused Facts
    causedFacts(N, Plan, E),
    lits_to_owl_strings(E, "", SFacts),
    % Reasons
    reasons(N, Plan, R),
    lits_to_owl_strings(R, "", SReasons),
    format(atom(S), "EquivalentClasses(:Action~w ObjectIntersectionOf(~w ObjectSomeValuesFrom(:inContext~w) ObjectSomeValuesFrom(:causes~w) ObjectSomeValuesFrom(:forReasons~w)))\n", [N, SName, SContexts, SFacts, SReasons]).

write_owl(Lines) :-
    open('./temp/temp.owl',write,OS),
    forall(member(M, Lines), writeln(OS,M)),
    write(OS, ")"),
    close(OS).

read_owl(O, Lines) :-
    open(O,read,Str), 
    read_owl_file(Str,Lines), 
    close(Str).
   
read_owl_file(Stream,[]):- 
    at_end_of_stream(Stream). 
read_owl_file(Stream,[X|L]):- 
    \+  at_end_of_stream(Stream), 
    read_line_to_codes(Stream,Codes),
    atom_chars(X, Codes),
    read_owl_file(Stream,L).

append_owl(Ontology, Additional, Onew) :-
    read_owl(Ontology, Lines),
    without_last(Lines, LinesWL),
    append(LinesWL, [Additional], Onew).