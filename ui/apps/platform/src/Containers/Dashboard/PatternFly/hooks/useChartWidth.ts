import useCurrentMediaBreakpoint, {
    PFBreakpoint,
} from 'hooks/patternfly/useCurrentMediaBreakpoint';

// These sizes are intended to maximize the size of the chart within a responsive
// container without distorting the underlying SVG element.
const chartWidths: Record<PFBreakpoint, number> = {
    /* Single column Grid layout for xs, sm, md */
    xs: 360,
    sm: 505,
    md: 700,
    /* Two column Grid layout for lg, xl, 2xl */
    lg: 430,
    xl: 520,
    '2xl': 655,
};

/**
 * Gets the recommended width in pixels for dashboard chart widgets. Used to
 * coordinate the size of the element display with the underlying SVG viewBox
 * for different screen sizes.
 */
export default function useChartWidth(): number {
    const currentBreakpoint = useCurrentMediaBreakpoint();
    return chartWidths[currentBreakpoint];
}
