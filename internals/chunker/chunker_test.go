package chunker

import (
	"testing"
)

func TestChunker(t *testing.T) {
	text := `The ocean looks calm today.

But underneath, there is constant motion. Currents shift slowly, sometimes unnoticed.

Some researchers say that the quietest systems hide the most complexity.



Data systems behave in a similar way.
At first glance, they look simple: input, processing, output.

But when the scale increases,
things start breaking in unexpected places.

Logs grow.

Queues fill up.
Workers retry.


And suddenly a tiny delay becomes a cascading failure.


This is why engineers spend so much time thinking about resilience.



Sometimes the solution is simple.

Add retries.

Add timeouts.

Add backoff.



But sometimes the solution is architectural.


Split services.

Introduce queues.

Cache aggressively.



And occasionally,
you discover the real issue wasn't infrastructure at all.

It was a tiny assumption hidden in a function written months ago.


One line.

One unchecked edge case.

One unexpected input.


That’s enough.



Debugging large systems often feels like searching for a needle in a haystack.

Except the haystack is constantly changing.


And sometimes the needle moves too.



So we build tools.

Observability systems.

Tracing.

Metrics.

Logs.


Because without visibility,

complex systems become impossible to reason about.



And when that happens,

even the smallest bug can look like chaos.`
	docID := "sywiq021wjsha_wu210352"
	userID := 1

	chunker := NewChunker(text, docID, int64(userID))
	chunks := chunker.Chunk()
	for i, c := range chunks {
		t.Logf("\n---- Chunk %d ----\n%v\n", i+1, c)
	}
}
