import { Component, OnInit } from '@angular/core';
import { Apollo } from 'apollo-angular';
import { gql } from 'apollo-angular';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html', 
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  searchQuery: string = '';
  states: any[] = [];
  map: google.maps.Map | undefined;  
  marker: google.maps.Marker | undefined;
  geocoder: google.maps.Geocoder | undefined;

  constructor(private apollo: Apollo) {}

  ngOnInit() {
    this.geocoder = new google.maps.Geocoder();
    this.map = new google.maps.Map(document.getElementById('googleMap') as HTMLElement, {
      zoom: 4,
      center: { lat: 37.0902, lng: -95.7129 }, // Default to the center of the USA
    });
    this.marker = new google.maps.Marker({
      map: this.map,
    });
  }

  onSearch() {
    if (this.searchQuery.length === 0) {
      this.states = [];
      return;
    }

    this.apollo
      .watchQuery({
        query: gql`
          query($filter: String!) {
            states(filter: $filter) {
              name
            }
          }
        `,
        variables: { filter: this.searchQuery },
      })
      .valueChanges.subscribe((result: any) => {
        this.states = result.data.states;
      });
  }

  selectState(state: any) {
    if (this.geocoder && this.map && this.marker) {
      this.geocoder.geocode({ address: state.name }, (results, status) => {
        if (status === 'OK' && results && results.length > 0) {
          const location = results[0].geometry.location;

          if (this.map && this.marker) {
            this.map.setCenter(location);
            this.map.setZoom(5);
            this.marker.setPosition(location);
            this.marker.setMap(this.map);
          }
        } else {
          console.error('Geocode failed for state: ' + state.name + ' with status: ' + status);
        }
      });
    }
  }
}
