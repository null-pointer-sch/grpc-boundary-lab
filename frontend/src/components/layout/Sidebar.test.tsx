import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Sidebar } from './Sidebar';

describe('Sidebar Component', () => {
    it('renders with navigation links', () => {
        render(<Sidebar collapsed={false} setCollapsed={vi.fn()} />);
        expect(screen.getByRole('menuitem', { name: /Dashboard/i })).toBeInTheDocument();
        expect(screen.getByRole('menuitem', { name: /Insights/i })).toBeInTheDocument();
    });

    it('triggers click events without crashing', () => {
        render(<Sidebar collapsed={false} setCollapsed={vi.fn()} />);
        const dashboardLink = screen.getByRole('menuitem', { name: /Dashboard/i });
        fireEvent.click(dashboardLink);
        // We just expect no crashing as links are mock # links
    });
});
