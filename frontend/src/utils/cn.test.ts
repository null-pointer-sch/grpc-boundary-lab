import { describe, it, expect } from 'vitest';
import { cn } from './cn';

describe('cn utility', () => {
    it('merges class names correctly', () => {
        expect(cn('px-2 py-1', 'bg-red-500')).toBe('px-2 py-1 bg-red-500');
    });

    it('handles conditional classes', () => {
        expect(cn('base-class', true && 'active', false && 'inactive')).toBe('base-class active');
    });

    it('resolves Tailwind conflicts using tailwind-merge', () => {
        // px-2 and px-4 conflict; px-4 should win
        expect(cn('px-2', 'px-4')).toBe('px-4');
        
        // bg-red-500 and bg-blue-500 conflict; bg-blue-500 should win
        expect(cn('bg-red-500 text-white', 'bg-blue-500')).toBe('text-white bg-blue-500');
    });

    it('handles arrays and complex objects', () => {
        expect(cn(['class1', 'class2'], { class3: true, class4: false })).toBe('class1 class2 class3');
    });
});
