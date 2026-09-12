# Accessibility

## Goal

Accessibility is a product requirement, not a final QA pass. The vertical slice should establish patterns that can scale to the full game.

## Minimum vertical-slice requirements

### Controls
- remappable gameplay controls,
- keyboard/mouse and controller parity for critical actions,
- avoid mandatory rapid repeated input where an alternative is practical,
- configurable steering/input sensitivity where supported.

### Visual
- UI scale options,
- readable contrast for critical text,
- critical state must not rely on color alone,
- subtitle/caption presentation with background/size options,
- clear focus/selection states.

### Motion
- reduced camera shake option,
- reduced motion/flash consideration,
- motion blur toggle where platform/engine supports it.

### Audio
- independent music/SFX/dialogue volume,
- captions/subtitles for critical voiced information,
- visual alternative for critical non-dialogue audio cues where practical.

## Racing-specific considerations

- checkpoint direction should use shape/icon/position as well as color,
- wrong-way and penalty warnings must be readable at speed,
- HUD density should have simplified options,
- assist settings must be separated from ranked integrity rules.

## Evidence

Accessibility readiness requires:
- keyboard-only/menu navigation review where applicable,
- controller navigation review,
- contrast/text-scale review,
- subtitle/caption review,
- documented known limitations.

No certification/conformance level is claimed until an explicit audit is completed.
