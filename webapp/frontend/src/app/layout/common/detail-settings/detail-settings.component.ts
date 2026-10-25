import {Component, Inject} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import {DetailService} from 'app/modules/detail/detail.service';

@Component({
    selector: 'app-detail-settings',
    templateUrl: './detail-settings.component.html',
    styleUrls: ['./detail-settings.component.scss']
})
export class DetailSettingsComponent {

    friendlyName: string;

    constructor(
        private _detailService: DetailService,
        private _dialogRef: MatDialogRef<DetailSettingsComponent>,
        @Inject(MAT_DIALOG_DATA) public data: { scrutiny_uuid: string, friendly_name: string }
    ) {
        this.friendlyName = data.friendly_name;
    }

    saveSettings(): void {
        this._detailService.setFriendlyName(this.data.scrutiny_uuid, this.friendlyName)
            .subscribe(friendlyName => this._dialogRef.close(friendlyName));
    }
}
