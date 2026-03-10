import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { ModeBadge } from './ModeBadge';

const mockMode = { protocol: 'grpc' as const, tls: false, tlsAvailable: true };

describe('ModeBadge Component', () => {
    it('renders grpc text', () => {
        render(<ModeBadge mode={mockMode} protocol="grpc" tlsEnabled={false} />);
        expect(screen.getByText('grpc')).toBeInTheDocument();
        expect(screen.queryByText('Secured')).not.toBeInTheDocument();
    });

    it('renders rest text', () => {
        render(<ModeBadge mode={mockMode} protocol="rest" tlsEnabled={false} />);
        expect(screen.getByText('rest')).toBeInTheDocument();
    });

    it('renders tls badge', () => {
        render(<ModeBadge mode={mockMode} protocol="grpc" tlsEnabled={true} />);
        expect(screen.getByText('Secured')).toBeInTheDocument();
    });
});
