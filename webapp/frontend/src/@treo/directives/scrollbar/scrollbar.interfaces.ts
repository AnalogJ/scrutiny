export class ScrollbarGeometry
{
    public x: number;
    public y: number;

    public w: number;
    public h: number;

    constructor(x: number, y: number, w: number, h: number)
    {
        this.x = x;
        this.y = y;
        this.w = w;
        this.h = h;
    }
}

export class ScrollbarPosition
{
    public x: number | 'start' | 'end';
    public y: number | 'start' | 'end';

    constructor(x: number | 'start' | 'end', y: number | 'start' | 'end')
    {
        this.x = x;
        this.y = y;
    }
}
