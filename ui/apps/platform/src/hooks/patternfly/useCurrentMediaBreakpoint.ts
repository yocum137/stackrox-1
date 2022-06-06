import { useState, useEffect } from 'react';
import uniqueId from 'lodash/uniqueId';

const rootStyle = getComputedStyle(document.documentElement);
const getPixelSize = (styleProp: string) => parseInt(styleProp.replace('px', ''), 10);
const pfXSmall = getPixelSize(rootStyle.getPropertyValue('--pf-global--breakpoint--xs'));
const pfSmall = getPixelSize(rootStyle.getPropertyValue('--pf-global--breakpoint--sm'));
const pfMedium = getPixelSize(rootStyle.getPropertyValue('--pf-global--breakpoint--md'));
const pfLarge = getPixelSize(rootStyle.getPropertyValue('--pf-global--breakpoint--lg'));
const pfXLarge = getPixelSize(rootStyle.getPropertyValue('--pf-global--breakpoint--xl'));
const pf2XLarge = getPixelSize(rootStyle.getPropertyValue('--pf-global--breakpoint--2xl'));

export const breakpoints = ['xs', 'sm', 'md', 'lg', 'xl', '2xl'] as const;

export type PFBreakpoint = typeof breakpoints[number];

const breakpointRanges: {
    name: PFBreakpoint;
    minWidth: number;
    maxWidth: number;
}[] = [
    { name: 'xs', minWidth: pfXSmall, maxWidth: pfSmall - 1 },
    { name: 'sm', minWidth: pfSmall, maxWidth: pfMedium - 1 },
    { name: 'md', minWidth: pfMedium, maxWidth: pfLarge - 1 },
    { name: 'lg', minWidth: pfLarge, maxWidth: pfXLarge - 1 },
    { name: 'xl', minWidth: pfXLarge, maxWidth: pf2XLarge - 1 },
    { name: '2xl', minWidth: pf2XLarge, maxWidth: 999999 },
];

let initialBreakpoint: PFBreakpoint = 'xs';
const listeners: Record<string, (breakpoint: PFBreakpoint) => void> = {};

breakpointRanges.forEach(({ name, minWidth, maxWidth }) => {
    const media = window.matchMedia(`(min-width: ${minWidth}px) and (max-width: ${maxWidth}px)`);
    if (media.matches) {
        initialBreakpoint = name;
    }
    const listener = (event) => {
        if (event.matches) {
            initialBreakpoint = name;
            Object.values(listeners).forEach((l) => l(name));
        }
    };
    media.addEventListener('change', listener);
});

/**
 * Hook that watches the current media breakpoint as defined by PatternFly CSS variables.
 */
export default function useCurrentMediaBreakpoint(): PFBreakpoint {
    const [id] = useState<string>(uniqueId('useCurrentBreakpoint'));
    const [breakpoint, setBreakpoint] = useState<PFBreakpoint>(initialBreakpoint);
    useEffect(() => {
        listeners[id] = setBreakpoint;

        return () => {
            delete listeners[id];
        };
    }, [id]);

    return breakpoint;
}
