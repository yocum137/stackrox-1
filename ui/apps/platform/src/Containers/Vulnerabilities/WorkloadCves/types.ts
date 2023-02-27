const detailsTabValues = ['Vulnerabilities', 'Resources'] as const;

export type DetailsTab = typeof detailsTabValues[number];

export function isDetailsTab(value: unknown): value is DetailsTab {
    return detailsTabValues.some((tab) => tab === value);
}
