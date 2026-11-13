# S<n>: <slice name>

Status: SPEC | BUILDING | VERIFYING | SHIPPED (date)
Work folder: `doc/<topic>_<YYYY-MM-DD>/`
Decisions: D-… · Assumptions: A-…

## Functions

| ID | Name | Signature | Units |
|---|---|---|---|
| `<PREFIX>-<TOPIC>-<NNN>` | `Name` | `func Name(…) (…, error)` | each input and the result |

## Per function

### `<ID>` `Name`
- **Source**: primary source, precise enough to find the line (edition, section, equation, table).
- **Best known**: time and space, and the lower bound with its source.
- **Approach**: the one chosen and why; **lost**: the alternative and why it lost.
- **Input modified**: yes or no.
- **Worked example** (first test): the inputs, the expected result, and where the expected result comes from.
- **Errors**: each sentinel and the input that produces it.

## Out of scope

What this slice does not do, and where it will live if it is ever done.
