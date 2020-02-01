package events

// ReactionHub is a hub for reactions. These are reactions all associated with a single entity.
// The reactions are all categorized and to be triggered on specific conditions.
type ReactionHub struct {
	OnCollision   []Reaction
	OnInteraction []Reaction
}

// NewReactionHub creates a reaction multiplexer with blank pipelines
func NewReactionHub() *ReactionHub {
	return &ReactionHub{[]Reaction{}, []Reaction{}}
}

// Condition under which a reaction should be triggered
const (
	ReactionOnCollision = iota
	ReactionOnInteraction
)

// Push - add a reaction to the ReactionHub based on the Condition T
func (r *ReactionHub) Push(T int, reaction Reaction) {
	switch T {
	case ReactionOnCollision:
		r.OnCollision = append(r.OnCollision, reaction)
		break
	case ReactionOnInteraction:
		r.OnInteraction = append(r.OnInteraction, reaction)
		break
	}
}

// HasReactions returns true if the reaction hub has reactions of the provided type.
func (r *ReactionHub) HasReactions(T int) bool {
	switch T {
	case ReactionOnCollision:
		return len(r.OnCollision) != 0
	case ReactionOnInteraction:
		return len(r.OnInteraction) != 0
	}
	return false
}
