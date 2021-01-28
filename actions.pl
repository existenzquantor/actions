#!/usr/bin/env swipl
:- initialization(main, main).
:- use_module("./core/core.pl", [classify_actions/1, classify_action/2, classify_plan/1, prepare_owl/2]).

main(Argv) :-
    nth0(0, Argv, Domain),
    nth0(1, Argv, Ontology),
    string_concat("./domains/", Domain, DomainPath0),
    string_concat(DomainPath0, ".pl", DomainPath),
    string_concat("./ontologies/", Ontology, OntologyPath0),
    string_concat(OntologyPath0, ".owl", OntologyPath),
    consult(DomainPath),
    plan(Plan),
    prepare_owl(Plan, OntologyPath),
    process(Argv).

process(Argv) :-
    length(Argv, 2),
    classify_actions(ActionTypes),
    classify_plan(PlanType),
    writeln([PlanType, ActionTypes]).
process(Argv) :-
    length(Argv, 3),
    nth0(2, Argv, plan),
    classify_plan(PlanType),
    writeln(PlanType).
process(Argv) :-
    length(Argv, 3),
    nth0(2, Argv, ArgActionNumber),
    atom_number(ArgActionNumber, ActionNumber),
    classify_action(ActionNumber,ActionTypes),
    writeln(ActionTypes).