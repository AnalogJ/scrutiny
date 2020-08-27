import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import {FileSizePipe} from "./file-size.pipe";
import { DeviceSortPipe } from './device-sort.pipe';

@NgModule({
    declarations: [
        FileSizePipe,
        DeviceSortPipe
    ],
    imports: [
        CommonModule,
        FormsModule,
        ReactiveFormsModule
    ],
    exports: [
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        FileSizePipe,
        DeviceSortPipe
    ]
})
export class SharedModule
{
}
