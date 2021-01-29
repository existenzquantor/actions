!/usr/bin/env swipl
:- initialization(main, main).
:- use_module("./core/core.pl", [classify_all_plans/3]).

main(Argv) :-
    nth0(0, Argv, Domain),
    string_concat("./domains/", Domain, DomainPath0),
    string_concat(DomainPath0, ".pl", DomainPath),
    consult(DomainPath),
    nth0(2, Argv, PlanLengthArg),
    atom_number(PlanLengthArg, PlanLength),
    nth0(1, Argv, Ontology),
    string_concat("./ontologies/", Ontology, OntologyPath0),
    string_concat(OntologyPath0, ".owl", OntologyPath),
    consult(DomainPath),
    classify_all_plans(PlanLength, OntologyPath, Result),
    writeln(Result).

