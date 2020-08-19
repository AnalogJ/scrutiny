import { animate, keyframes, style, transition, trigger } from '@angular/animations';

// -----------------------------------------------------------------------------------------------------
// @ Shake
// -----------------------------------------------------------------------------------------------------
const shake = trigger('shake',
    [

        // Prevent the transition if the state is false
        transition('void => false', []),

        // Transition
        transition('void => *, * => true',
            [
                animate('{{timings}}',
                    keyframes([
                        style({
                            transform: 'translate3d(0, 0, 0)',
                            offset   : 0
                        }),
                        style({
                            transform: 'translate3d(-10px, 0, 0)',
                            offset   : 0.1
                        }),
                        style({
                            transform: 'translate3d(10px, 0, 0)',
                            offset   : 0.2
                        }),
                        style({
                            transform: 'translate3d(-10px, 0, 0)',
                            offset   : 0.3
                        }),
                        style({
                            transform: 'translate3d(10px, 0, 0)',
                            offset   : 0.4
                        }),
                        style({
                            transform: 'translate3d(-10px, 0, 0)',
                            offset   : 0.5
                        }),
                        style({
                            transform: 'translate3d(10px, 0, 0)',
                            offset   : 0.6
                        }),
                        style({
                            transform: 'translate3d(-10px, 0, 0)',
                            offset   : 0.7
                        }),
                        style({
                            transform: 'translate3d(10px, 0, 0)',
                            offset   : 0.8
                        }),
                        style({
                            transform: 'translate3d(-10px, 0, 0)',
                            offset   : 0.9
                        }),
                        style({
                            transform: 'translate3d(0, 0, 0)',
                            offset   : 1
                        })
                    ])
                )
            ],
            {
                params: {
                    timings: '0.8s cubic-bezier(0.455, 0.03, 0.515, 0.955)'
                }
            }
        )
    ]
);

export { shake };
