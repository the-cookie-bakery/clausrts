# Clausewitz Reference — CLAUSRTS

Source text: `references/vomkriege_book/_VKwholetext.htm` — the complete 1832
German *Vom Kriege* (clausewitzstudies.org). Each game concept is mapped to the
book chapter that defines it.

## Trinity — *wunderliche Dreifaltigkeit* (Buch I, Kap. 1: *Was ist der Krieg?*)

> „...ihm herrschenden Tendenzen eine wunderliche Dreifaltigkeit, zusammengesetzt
> aus der ursprünglichen Gewaltsamkeit seines Elementes, dem Haß und der
> Feindschaft, die wie ein blinder Naturtrieb anzusehen sind, aus dem Spiel der
> Wahrscheinlichkeiten und des Zufalls, die ihn zu einer freien Seelentätigkeit
> machen, und aus der untergeordneten Natur eines politischen Werkzeuges, wodurch
> er dem bloßen Verstande anheimfällt."

- Passion / hatred → **People (Volk)**
- Chance / probability → **Military (Heer)**
- Reason / policy → **Government (Regierung)**

## War as instrument of policy (Buch I, Kap. 6B / Kap. 8)

> „Der Krieg ist eine bloße Fortsetzung der Politik mit anderen Mitteln."

The Government bar encodes this: if military means drift from the political aim,
legitimacy collapses.

## Friction — *Friktion im Kriege* (Buch I, Kap. 7)

> „Die Gefahr im Kriege gehört zur Friktion desselben ..."

The `Friction` value on every formation. Raised by fog, distance from HQ, ECM,
fatigue, terrain; lowered by reconnaissance and good command.

## Fog / intelligence — *Nachrichten im Kriege* (Buch I, Kap. 6)

> „...aus den Nachrichten im Kriege hervorgeht ... die meisten davon sind falsch."

Per-faction visibility grid. GPS/drones/satellites dissolve it; scramblers
restore it.

## Center of gravity — *Schwerpunkt* (Buch VIII: *Kriegsplan*; Buch VI)

> „...daß man also die strategischen Wirkungen auf gewisse Schwerpunkte
> zurückführen kann ..."

The `Gravity` value and Schwerpunkt detection. Strike the enemy center of
gravity; protect your own.

## Culmination — *Kulminationspunkt des Angriffs* (Buch VII, Kap. 5)

> „Über den Kulminationspunkt des Sieges ..."

The `Culmination` value: offensive power decays with distance/time; past the
peak, defense is stronger.

## Economy of force — *Ökonomie der Kräfte* / *Sammlung der Kräfte im Raum* (Buch III)

The `Focus` value: concentration in space/time.

## Reciprocal action — *Wechselwirkung von Angriff und Verteidigung* (Buch I)

The counter matrix: every weapon has a counter; no universal means.

## Volksbewaffnung (Buch VI, Kap. 26)

High People bar unlocks partisan/militia units in home territory.

## Spannung und Ruhe (Buch III, Kap. 18)

The pause/plan mode: plan while paused, resolve in real time.
