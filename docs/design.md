# Design — CLAUSRTS

## 1. Vision

A 2D, tile-based RTS (real-time with pause) for Linux/ebitengine that teaches
Clausewitz's *Vom Kriege* through a modern hybrid-war setting. The player does
not "gather resources": war is *not* an independent economic act, it is a
continuation of politics by other means. Instead, the player manages the three
elements of Clausewitz's **Trinity** while fighting a war whose outcome is
determined by **friction, fog, focus, gravity (Schwerpunkt), and culmination**.

## 2. Learning objectives (in teaching order)

1. **War is an instrument of policy.** Means must stay aligned with the political
   aim, or the war loses legitimacy and support (Buch I, Kap. 1 & 6B).
2. **"Everything in war is very simple, but the simplest thing is difficult"**
   — friction distinguishes real war from war on paper (Buch I, Kap. 7).
3. **Fog / imperfect intelligence.** Reports are unreliable; reconnaissance
   dissolves fog, jamming restores it (Buch I, Kap. 6).
4. **The culminating point.** The offensive decays with distance/time; a
   commander must recognize the peak and consolidate (Buch VII, Kap. 4–5).
5. **Center of gravity (Schwerpunkt).** Concentrate against the enemy's
   center of gravity and protect your own (Buch VIII).
6. **Wechselwirkung (reciprocal action).** Every weapon has a counter; there is
   no universal means (Buch I).

## 3. Core gameplay loop

1. **Plan (pause).** Read fog, issue orders, merge/split units between
   formations, attach support assets (scrambler + tanks = combined arms).
2. **Resolve (real time).** Movement, combat, friction, morale and supply tick;
   the three Trinity bars drift.
3. **Reflect.** The teaching layer surfaces the relevant *Vom Kriege* quote when
   a concept becomes active (e.g. friction spikes → "Friktion im Kriege, Kap. 7").

There is **no base-building**. Forces are pre-set per scenario; reinforcements
trickle in from the map edge, gated by the People + Government bars.

## 4. The Trinity — three support bars

From Buch I, Kap. 1 (*"Was ist der Krieg?"*):

> „...eine wunderliche Dreifaltigkeit, zusammengesetzt aus der ursprünglichen
> Gewaltsamkeit seines Elementes, dem Haß und der Feindschaft, ... aus dem Spiel
> der Wahrscheinlichkeiten und des Zufalls, ... und aus der untergeordneten Natur
> eines politischen Werkzeuges ..."

| Bar | Clausewitz element | Game meaning | Feeds | Depletes when |
|---|---|---|---|---|
| **People (Volk)** | Passion / hatred | National morale, war enthusiasm | Reinforcement rate, partisan/militia units (Volksbewaffnung) | Casualties, long war, civilian loss |
| **Military (Heer)** | Chance / genius | Operational strength & readiness | Focus (better command = lower friction) | Attrition, over-extension |
| **Government (Regierung)** | Reason / policy | Political capital & legitimacy | Command authority, cyber/strategic assets, negotiation | War drifting from the political objective |

The bars are **interlinked** (the Trinity): actions trade between them. If any
bar reaches zero the game ends (legitimacy collapse / army collapse / revolution).

## 5. Clausewitz properties on formations

Every formation carries a `Dynamics` struct that encodes the concepts the game
teaches. These are the values the player can observe and influence:

| Property | Meaning | Raised by | Lowered by |
|---|---|---|---|
| `Friction` | "War on paper vs real war" | Fog, distance from HQ, ECM/jamming, fatigue, terrain | Recon (GPS/drones), good command, roads |
| `Focus` | Concentration in space/time (Ökonomie der Kräfte) | Concentrating units, GPS targeting | Dispersion, enemy scrambling |
| `Gravity` | Schwerpunkt weight of a location/force | Concentrating forces, HQs, objectives | Destruction of that force |
| `Fog` | Unknown enemy positions (Nachrichten) | GPS scrambling, concealment | Recon, drones, satellites |
| `Culmination` | Progress toward the point where attack flips to disadvantage | Advancing far from supply, high friction | Consolidating, resupply, defense |
| `Moral` | Cohesion / moralische Hauptpotenzen | Defense of homeland, success | Friction, casualties |
| `Supply` | Lines of communication (Verbindungslinien) | Proximity to depot | Long advances, interdiction |

These interact: high friction reduces focus, which reduces attack power, moving
the culmination point closer. The teaching layer explains *why* a value changed.

## 6. Modern hybrid-war units & counter matrix

Each sub-unit type is a `entities/Unit` carrying static stats plus
Clausewitz-property modifiers. The matrix is a closed loop — teaching
*Wechselwirkung* (every weapon has a counter).

| Asset | Effect | Hard counter |
|---|---|---|
| Infantry (mech) | Holds ground, cheap, garrisons | Artillery, armor |
| Armor / Tanks | Breakthrough (focus) | ATGM infantry, drones, AT mines; **long advance → culmination/friction** |
| Rockets / precision artillery | Ranged attrition, breaks cohesion | **Bunker** (shield), air defense (intercept), dispersion, counter-battery |
| Drones (recon/strike) | Dissolve fog, reduce friction, light strike | **Air defense**, EW jamming, anti-drone |
| GPS / satellite support | Dissolve fog, reduce friction, raise focus | **GPS scrambler** (regional) |
| GPS scrambler / EW | Raises enemy friction + fog, hides movement | Direction-finding (EW hunt), hardened comms, inertial nav (partial) |
| Cyber unit | Raises enemy friction, disables buildings, steals intel | **Firewall** (cyber defense) |
| Air defense | Counters drones/rockets | Artillery (SEAD), cyber, ground attack |
| Engineer | Builds bunker/firewall/trenches | Artillery suppression |

### Implemented unit stats

| Kind (`entities/<name>.go`) | HP | Attack | Range | Speed | Counters (2×) |
|---|---|---|---|---|---|
| infantry | 80 | 12 | 1 | 3 | airdefense, cyber |
| tank | 220 | 45 | 2.5 | 4 | infantry, artillery, airdefense, cyber |
| artillery | 90 | 60 | 8 | 2 | infantry, tank, engineer |
| drone *(flies)* | 40 | 18 | 2 | 6 | tank, artillery, infantry |
| airdefense | 90 | 35 | 4 | 3 | drone, artillery |
| cyber | 50 | 25 | 1.5 | 3 | gps, scrambler, radar, depot |
| scrambler | 60 | 5 | 1 | 3 | drone, gps |
| gps | 50 | 5 | 1 | 3 | — |
| engineer | 70 | 8 | 1 | 3 | — |

### Buildings (homebase = your Schwerpunkt cluster)

| Building | Role | Implemented |
|---|---|---|
| **HQ** | Center of gravity; losing it is catastrophic | ✅ |
| **Supply depot** | Defines the culmination curve (Verbindungslinien) | ✅ |
| **Radar / comms** | Reduces fog & friction (command radius) | ✅ |
| **Bunker** | Shields vs rockets/artillery (currently shields its own tile) | ✅ |
| **Firewall** | Cyber defense | ✅ |
| **Power grid** | Cyber target | ❌ (planned) |

The implemented combat rule is simple: an attacker deals **2× damage** to unit or
building kinds listed in its `Counters` set. Support units (GPS, scrambler,
drone) also emit **recon/jam auras** that raise/lower nearby fog and friction via
`systems/dynamics.go`. Drones can **fly over water** (`UnitKind.CanFly`).

## 7. Per-unit control & formations

- A **Formation** (battalion) = leader token + list of sub-units (1×1 tile each).
- Sub-units can be **transferred** between formations: reinforce a depleted
  battalion, or attach a scrambler/air-defense to a tank battalion (combined
  arms lowers friction, raises focus, but the support asset becomes a target).
- Formation `Dynamics` are **computed from composition** — this is where
  "different parts of the military interact" to raise/lower Clausewitz values.

**Status:** the `Formation` model exists and each unit records its formation, but
the merge/split transfer UI is **not yet wired** — units currently act
independently. This is part of M4/M5.

## 8. Scenarios

Each scenario teaches one or two concepts (see [roadmap.md](roadmap.md)):

1. **"Was ist der Krieg?"** — tutorial.
2. **"Friktion im Kriege"** — cross a jammed map.
3. **"Kulminationspunkt"** — recognize the offensive's peak.
4. **"Schwerpunkt"** — strike/protect the center of gravity.
5. **"Ökonomie der Kräfte"** — concentrate vs disperse; full counter-matrix.
6. **"Kriegsplan"** — full campaign to a political (bounded) end.
