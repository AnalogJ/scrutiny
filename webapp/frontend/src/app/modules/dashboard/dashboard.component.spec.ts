import {tooltipViewportOffset} from './dashboard.component';

describe('tooltipViewportOffset', () => {
    const viewport = {width: 1000, height: 800};
    const box = (left: number, top: number, width: number, height: number) =>
        ({left, top, right: left + width, bottom: top + height});

    it('leaves a tooltip that already fits where it is', () => {
        expect(tooltipViewportOffset(box(100, 100, 200, 300), viewport)).toEqual({dx: 0, dy: 0});
    });

    it('lifts a tooltip that runs off the bottom', () => {
        expect(tooltipViewportOffset(box(100, 600, 200, 300), viewport)).toEqual({dx: 0, dy: -108});
    });

    it('pulls a tooltip that runs off the right back to the left', () => {
        expect(tooltipViewportOffset(box(900, 100, 200, 100), viewport)).toEqual({dx: -108, dy: 0});
    });

    it('pins a tooltip taller than the viewport to the top rather than lifting it out of view', () => {
        expect(tooltipViewportOffset(box(100, 200, 200, 900), viewport)).toEqual({dx: 0, dy: -192});
    });

    it('does not push a tooltip off the left while pulling it in from the right', () => {
        expect(tooltipViewportOffset(box(-50, 100, 1200, 100), viewport)).toEqual({dx: 58, dy: 0});
    });

    it('settles: re-applying the offset to the moved tooltip is a no-op', () => {
        const start = box(900, 600, 200, 300);
        const {dx, dy} = tooltipViewportOffset(start, viewport);
        const moved = box(start.left + dx, start.top + dy, 200, 300);
        expect(tooltipViewportOffset(moved, viewport)).toEqual({dx: 0, dy: 0});
    });
});
