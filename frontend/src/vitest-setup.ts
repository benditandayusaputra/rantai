import '@testing-library/jest-dom/vitest';
import { vi } from 'vitest';

vi.mock('$env/dynamic/public', () => ({ env: {} }));

HTMLDialogElement.prototype.showModal = function (this: HTMLDialogElement) {
	this.open = true;
};
HTMLDialogElement.prototype.close = function (this: HTMLDialogElement) {
	this.open = false;
	this.dispatchEvent(new Event('close'));
};
