# actions

This is designed as a tool for causal and terminological reasoning about actions in action plans. That is, given an action-type ontology, an action-domain specification, and an action plan, the tool will determine under which actions types each of the actions in the action plan falls.

The tool is originally motivated by the desire to modeling the phenomena that an action can be described in very different terms. For instance, consider the Trolley Problem (see https://en.wikipedia.org/wiki/Trolley_problem). The act of pulling the lever can be described as just pulling a lever, or as an act of killing one person, or as an act of rescuing the lives of five persons. Depending on the description chosen, the action be evaluated differently with respect to moral permissibility. 

# Installation

To run the tool, SWI Prolog (https://www.swi-prolog.org/) must be installed on the system. SWI Prolog is available for many operating systems. In a MacOS environment, you can install SWI Prolog using brew by typing ```brew install swipl```, under Ubuntu Linux it is ```apt install swipl```, and binaries for Windows are also available on the SWI Prolog download website. 

To get started, first clone this github repository. Currently, it is assumed that the folder containing the causality library (https://github.com/existenzquantor/causality) is located next to the hera2 folder, i.e., in some common parent folder. Thus, you may need to clone the causality github repository, as well.

# Tutorial

## Modeling an Action-Type Ontology

To compute the action types of actions in an action plan, the tool requires an action-type ontology. You are pretty free to define an action-type ontology you like. One convenient way of doing that is using Protégé (see https://protege.stanford.edu). There are, however, two constraints to consider:

* The Ontology has to be stored in the ontologies folder, it has to have the ending ".owl", and it must be written in Functional OWL Syntax. Protégé can import and expert in the format.
* Every action in the action plan will be turned into a particular representation consisting of four components:
  * the action's given name
  * the context in which the action is performed (i.e., the state it is applied to during execution),
  * the facts the action causes,
  * the reasons the action is performed for.

For instance, considering the Trolley Problem again, the action "pull" occurring in an action plan, may be represented like this:
```
EquivalentClasses(:Action0 
  ObjectIntersectionOf(:pull 
   ObjectSomeValuesFrom(:inContext ObjectIntersectionOf( ObjectComplementOf(:dead5) ObjectComplementOf(:dead1) ObjectComplementOf(:left))) 
   ObjectSomeValuesFrom(:causes ObjectIntersectionOf( :dead1 :left ObjectComplementOf(:dead5))) 
   ObjectSomeValuesFrom(:forReasons ObjectIntersectionOf( ObjectComplementOf(:dead5)))
  ))
```

This line says that Action0 (thus, the first action in the plan) has "pull" as given name and it is is performed in a state where the five persons on the one track are not dead, the other person on the other track is not dead, the trolley is not directed towards the left track. Moreover, the action causes the death of the one person, it causes the trolley to be directed onto the left track, and it causes the survival of the five persons. Finally, the action is performed for the reason to save the lives of the five persons.

Hence, when you design your action-type ontology, you may want to make use of the relations "inContext", "causes", and "forReasons" as indicated by this example.


## Modeling an Action Domain

## Running the Tool

