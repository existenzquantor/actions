# actions

This is designed as a tool for causal and terminological reasoning about actions in action plans. That is, given an action-type ontology, an action-domain specification, and an action plan, the tool will determine under which actions types each of the actions in the action plan falls.

The tool is originally motivated by the desire to modeling the phenomena that an action can be described in very different terms. For instance, consider the Trolley Problem (see https://en.wikipedia.org/wiki/Trolley_problem). The act of pulling the lever can be described as just pulling a lever, or as an act of killing one person, or as an act of rescuing the lives of five persons. Depending on the description chosen, the action be evaluated differently with respect to moral permissibility. 

# Installation

To run the tool, SWI Prolog (https://www.swi-prolog.org/) must be installed on the system. SWI Prolog is available for many operating systems. In a MacOS environment, you can install SWI Prolog using brew by typing ```brew install swipl```, under Ubuntu Linux it is ```apt install swipl```, and binaries for Windows are also available on the SWI Prolog download website. 

To get started, first clone this github repository. Currently, it is assumed that the folder containing the causality library (https://github.com/existenzquantor/causality) is located next to the hera2 folder, i.e., in some common parent folder. Thus, you may need to clone the causality github repository, as well.

# Tutorial

## Modeling an Action-Type Ontology

## Modeling an Action Domain

## Running the Tool

