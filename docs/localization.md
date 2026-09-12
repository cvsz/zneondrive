# Localization

## Initial language strategy

The project should be authored localization-first even if early prototypes ship in one language.

Recommended initial commercial set:
- English,
- Thai,
- Simplified Chinese,
- Japanese,
- Korean.

Final launch languages depend on commercial/platform decisions.

## Content rules

- stable IDs never contain translated display strings,
- dialogue/UI text should use localization keys,
- avoid concatenating grammar-sensitive sentence fragments,
- UI layouts must tolerate expansion,
- dates/numbers/currency use locale-aware formatting where applicable,
- player-generated text is not translated by default.

## Narrative workflow

1. Lock source text/version.
2. Export localization units with context.
3. Translate.
4. Linguistic review.
5. Import.
6. In-game layout/cinematic review.
7. Fix source/context issues rather than hardcoding per-language exceptions where possible.

## Vehicle/technical terminology

Maintain a termbase for:
- parts,
- tuning concepts,
- race disciplines,
- factions,
- district names,
- DRIVE ZERO terminology.

Proper nouns require an explicit translate/transliterate/keep policy.

## Evidence

Localization completeness should be tracked by key coverage and in-context review, not only by translated-file percentage.
