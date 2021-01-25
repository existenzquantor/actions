#!/usr/bin/env swipl
:- initialization(main, main).
:- use_module("./core/core.pl", [classify_actions/1, classify_plan/1, prepare_owl/1]).

main(Argv) :-
    nth0(0, Argv, Domain),
    nth0(1, Argv, Ontology),
    string_concat("./domains/", Domain, DomainPath0),
    string_concat(DomainPath0, ".pl", DomainPath),
    string_concat("./ontologies/", Ontology, OntologyPath0),
    string_concat(OntologyPath0, ".owl", OntologyPath),
    consult(DomainPath),
    prepare_owl(OntologyPath),
    classify_actions(ActionTypes),
    classify_plan(PlanType),
    writeln([PlanType, ActionTypes]).