package monthly

func (r *Registry) ContainsHash(checksum string) bool {
	_, exists := r.checksumIndex[checksum]
	return exists
}
