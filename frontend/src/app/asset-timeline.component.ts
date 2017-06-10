import { Component, Input, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Http } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/switchMap';

import { AssetService } from './asset.service';

@Component({
    selector: 'app-root',
    templateUrl: './asset-timeline.component.html',
    providers: [AssetService],
})
export class AssetTimelineComponent implements OnInit {

    constructor(
        private assetService: AssetService,
        private route: ActivatedRoute) {
    }

    ngOnInit(): void {
    }
}

