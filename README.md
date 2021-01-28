# actions

This is designed as a tool for causal and terminological reasoning about actions in action plans. That is, given an action-type ontology, an action-domain specification, and an action plan, the tool will determine under which actions types each of the actions in the action plan falls.

The tool is originally motivated by the desire to modeling the phenomena that an action can be described in very different terms. For instance, consider the Trolley Problem (see https://en.wikipedia.org/wiki/Trolley_problem). The act of pulling the lever can be described as just pulling a lever, or as an act of killing one person, or as an act of rescuing the lives of five persons. Depending on the description chosen, the action be evaluated differently with respect to moral permissibility. 

The approach was first presented in: 
*Lindner F. Permissibility-under-a-description reasoning for deontological robots. 
To Appear in: Proceedings of RoboPhilosophy 2020; 2020*. 

A preprint is available here: http://existenzquantor.info/wp-content/uploads/2020/11/FL_RP_2020_Preprint.pdf


# Installation

To run the tool, SWI Prolog (https://www.swi-prolog.org/) must be installed on the system. SWI Prolog is available for many operating systems. In a MacOS environment, you can install SWI Prolog using brew by typing ```brew install swipl```, under Ubuntu Linux it is ```apt install swipl```, and binaries for Windows are also available on the SWI Prolog download website. Moreover, a Java Virtual Machine must be installed on the system.

To get started, first clone this github repository. Currently, it is assumed that the folder containing the causality library (https://github.com/existenzquantor/causality) is located next to the hera2 folder, i.e., in some common parent folder. Thus, you may need to clone the causality github repository, as well.

# Tutorial

## Modeling an Action-Type Ontology

To compute the action types of actions in an action plan, the tool requires an action-type ontology. You are pretty free to define an action-type ontology you like. One convenient way of doing that is using Protégé (see https://protege.stanford.edu). There are, however, two constraints to consider:

* The Ontology has to be stored in the ontologies folder, it has to have the ending ".owl", and it must be written in Functional OWL Syntax. Protégé can import and export that format.
* Every action in the action plan will be turned into a particular representation consisting of four components:
  * the action's given name
  * the context in which the action is performed (i.e., the state it is applied to during execution),
  * the facts the action causes,
  * the reasons the action is performed for.

For instance, considering the Trolley Problem again, the action "pull" occurring in an action plan, may be represented like this:
```
EquivalentClasses(:Action0 
  ObjectIntersectionOf(:pull 
   ObjectSomeValuesFrom(action:inContext ObjectIntersectionOf(ObjectComplementOf(:dead5) ObjectComplementOf(:dead1) ObjectComplementOf(:left))) 
   ObjectSomeValuesFrom(action:causes ObjectIntersectionOf(:dead1 :left ObjectComplementOf(:dead5))) 
   ObjectSomeValuesFrom(action:forReasons ObjectIntersectionOf(ObjectComplementOf(:dead5)))
  ))
```

This line says that Action0 (thus, the first action in the plan) has "pull" as given name and it is is performed in a state where the five persons on the one track are not dead, the other person on the other track is not dead, the trolley is not directed towards the left track. Moreover, the action causes the death of the one person, it causes the trolley to be directed onto the left track, and it causes the survival of the five persons. Finally, the action is performed for the reason to save the lives of the five persons.

Hence, when you design your action-type ontology, you may want to make use of the relations "inContext", "causes", and "forReasons" as indicated by this example.

An example ontology for the Trolley Problem can be found here: https://github.com/existenzquantor/actions/blob/master/ontologies/Trolley.owl This ontology defines the action types Kill (causing death) and Rescue (causing non-death).


## Modeling an Action Domain

The action-type ontology is complemented by an action-domain model. Action-domain models live in the "domains" directory. These model the actions in terms of preconditions and effects, they model events, the agent's goals, the initial state, and the plan the agent actually performs. As an example, consider the trolley domain here: https://github.com/existenzquantor/actions/blob/master/domains/trolley.pl . It says that pulling the lever results in the trolley be directed onto the left track, if it is not directed onto the left track before the action is performed. Thus, the first argument of the effect fact is the action name, the second one is a list of literals that must hold as precondition for the effect to obtain in the successor state (third argument, again a list of literals).

By convention, effect facts that have numbers as first arguments are external events. That is, if the event is called n, then after the n-th action, the event will fire, that is, the event's effects will obtain if the event's preconditions hold after the n-th action. n starts counting with 0. Thus, in the particular example, the one person will be dead if pull is the first action in the plan.

Note that literals used in the action domain will be represented as concepts in the action-type representation of the actions (see above).


## Running the Tool

Finally, to retrieve the action descriptions, the following command must be executed in the shell: ```swipl actions.pl <domain> <ontology>```. In the example, we type ```swipl actions.pl trolley Trolley```. The tool will find ```trolley.pl``` in the domain directory, and ```Trolley.owl``` in the ontology directory. As a result, we get ```[[Rescue, Kill, plan], [[Rescue,Kill,pull]]]```. The first entry is a list of action types denoting the descriptions for the entire plan. The second entry is a list of lists of action types: one list for each actions in the action plan. In the Trolley example, there is only one action (viz., pull), and it can be described as from type Rescue, Kill, or just pull.

Also try the other example domain:
```bash
> swipl actions.pl flipswitch2 FlipSwitch
[[plan], [[FlipSwitch,TurnOn],[FlipSwitch,TurnOff]]]
```
In this example, the action FlipSwitch is performed twice (starting from a state where the light is off). The plan as a whole does not cause anything, thus, it has only a trivial description: It's a plan. The agent of that plan may explain its behavior by saying something like: _I switch the light on and off because I want to_. However, the individual actions do have richer descriptions: the first FlipSwitch action is also of type turning-the-light-on (TurnOn), and the second FlipSwitch action turns the light off (TurnOff).

The file ```actions.pl``` can also be made executable (```chmod +x actions.pl```). Then you can also run:
```bash
> ./actions.pl flipswitch2 FlipSwitch
[[plan], [[FlipSwitch,TurnOn],[FlipSwitch,TurnOff]]]
```

It is also possible to explicitly query the tool for the set of descriptions of the plan or for some particular action within that plan. For querying the descriptions of the plan, just append "plan" as third argument:
```bash
> ./actions.pl flipswitch2 FlipSwitch plan
[plan]
```

For querying the n-th action in the plan, just append "n" as third argument:
```bash
> ./actions.pl flipswitch2 FlipSwitch 0
[FlipSwitch,TurnOff]
> ./actions.pl flipswitch2 FlipSwitch 1
[FlipSwitch,TurnOn]
```
Note that the count starts with 0, that is, the first action in the plan has index 0. Generally, querying for a single action is faster than querying the whole information using the two-arguments command and then extracting the desired information. However, querying multiple actions action-by-action may be much slower than querying the whole bunch of information and then extracting the desired information from the output.
