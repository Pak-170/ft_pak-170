import { writable } from 'svelte/store';

export const loginHasError = writable(false);
export const passwordMismatch = writable(false);
export const loginLoading = writable(false);
export const loginSuccess = writable(false);
