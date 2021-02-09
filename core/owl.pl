- module(actions_owl, [prepare_owl/2]).
:- use_module("core.pl", [names/2, contexts/4, causedFacts/3, causedFacts/1, reasons/3, reasons/2]).
:- use_module("helpers.pl", [without_last/2]).

prepare_owl(Plan, OntologyPath) :-
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
    format(atom(SF), " ObjectIntersectionOf(~w owl:Thing owl:Thing)", S).
lits_to_owl_strings([A|R], S, F) :-
    lits_to_owl_strings(A, S, F2),
    lits_to_owl_strings(R, F2, F).

lits_to_owl_strings_or(not(A), S, F) :-
    format(atom(X), " ObjectComplementOf(:~w)", A),
    string_concat(S, X, F),!.
lits_to_owl_strings_or(A, S, F) :-
    \+is_list(A),
    format(atom(X), " :~w", A),
    string_concat(S, X, F),!.
lits_to_owl_strings_or([], S, SF) :-
    format(atom(SF), " ObjectUnionOf(~w owl:Thing owl:Thing)", S).
lits_to_owl_strings_or([A|R], S, F) :-
    lits_to_owl_strings(A, S, F2),
    lits_to_owl_strings(R, F2, F).

to_owl(Plan, S) :-
    names(Plan, Names),
    length(Names, N),
    to_owl(Plan, 0, N, "", S1),
    to_owl(S2),
    string_concat(S1, S2, S).
to_owl(_, N, N, S, S).
to_owl(Plan, N, Nmax, Str, S) :-
    N < Nmax,
    to_owl(Plan, N, O),
    string_concat(Str, O, Str2),
    N2 is N + 1,
    to_owl(Plan, N2, Nmax, Str2, S).
to_owl(Plan, N, S) :-
    N > -1,
     % Names
    names(Plan, Names),
    nth0(N, Names, Action),
    lits_to_owl_strings(Action, "", SName),
    % Context
    init(C),
    contexts(Plan, C, [], Contexts),
    nth0(N, Contexts, Context),
    lits_to_owl_strings(Context, "", SContexts),
    lits_to_owl_strings_or(Context, "", SContextsOr),
    % Caused Facts
    causedFacts(N, Plan, E),
    lits_to_owl_strings(E, "", SFacts),
    lits_to_owl_strings_or(E, "", SFactsOr),
    % Reasons
    reasons(N, Plan, R),
    lits_to_owl_strings(R, "", SReasons),
    lits_to_owl_strings_or(R, "", SReasonsOr),
    format(atom(S), "EquivalentClasses(:Action~w ObjectIntersectionOf(~w ObjectSomeValuesFrom(action:inContext~w) ObjectAllValuesFrom(action:inContext~w) ObjectSomeValuesFrom(action:causes~w) ObjectAllValuesFrom(action:causes~w) ObjectSomeValuesFrom(action:forReasons~w) ObjectAllValuesFrom(action:forReasons~w)))\n", [N, SName, SContexts, SContextsOr, SFacts, SFactsOr, SReasons, SReasonsOr]).
to_owl(S) :-
    init(IS),
    lits_to_owl_strings(IS, "", SContexts),
    lits_to_owl_strings_or(IS, "", SContextsOr),
    % Caused Facts
    causedFacts(E),
    lits_to_owl_strings(E, "", SFacts),
    lits_to_owl_strings_or(E, "", SFactsOr),
    % Reasons
    reasons(E, R),
    lits_to_owl_strings(R, "", SReasons),
    lits_to_owl_strings_or(R, "", SReasonsOr),
    format(atom(S), "EquivalentClasses(:Plan0 ObjectIntersectionOf(:Plan ObjectSomeValuesFrom(action:inContext~w) ObjectAllValuesFrom(action:inContext~w) ObjectSomeValuesFrom(action:causes~w) ObjectAllValuesFrom(action:causes~w) ObjectSomeValuesFrom(action:forReasons~w) ObjectAllValuesFrom(action:forReasons~w)))\n", [SContexts, SContextsOr, SFacts, SFactsOr, SReasons, SReasonsOr]).


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
    (\+ (Codes = [35 | _]) ->
    atom_chars(X, Codes);
    atom_chars(X, [])),
    read_owl_file(Stream,L).

append_owl(Ontology, Additional, Onew) :-
    read_owl(Ontology, Lines),
    without_last(Lines, LinesWL),
    append(LinesWL, [Additional], Onew).