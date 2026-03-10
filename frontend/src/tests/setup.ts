import '@testing-library/jest-dom';
import { cleanup } from '@testing-library/react';
import { afterEach } from 'vitest';

afterEach(() => {
  cleanup();
});

class ResizeObserver {
  observe(): void {
    // Stub for testing
  }

  unobserve(): void {
    // Stub for testing
  }

  disconnect(): void {
    // Stub for testing
  }
}
;(globalThis as any).ResizeObserver = ResizeObserver;
