import { writable } from 'svelte/store';

export const loginHasError = writable(false);
export const passwordMismatch = writable(false);
