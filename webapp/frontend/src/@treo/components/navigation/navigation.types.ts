export interface TreoNavigationItem
{
    id?: string;
    title?: string;
    subtitle?: string;
    type: 'aside' | 'basic' | 'collapsable' | 'divider' | 'group' | 'spacer';
    hidden?: (item: TreoNavigationItem) => boolean;
    disabled?: boolean;
    link?: string;
    externalLink?: boolean;
    exactMatch?: boolean;
    function?: (item: TreoNavigationItem) => void;
    classes?: string;
    icon?: string;
    iconClasses?: string;
    badge?: {
        title?: string;
        style?: 'rectangle' | 'rounded' | 'simple',
        background?: string;
        color?: string;
    };
    children?: TreoNavigationItem[];
    meta?: any;
}

export type TreoVerticalNavigationAppearance = string;
export type TreoVerticalNavigationMode = 'over' | 'side';
export type TreoVerticalNavigationPosition = 'left' | 'right';
