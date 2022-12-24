---
title: Location
layout: default
grand_parent: Documentation
parent: Actions
nav_order: 4
---

# Location Actions

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

### `getCurrentLocation()`

Get the users current location.

---

### `getLocationDetail(location,detail)`

Get `detail` of `location`.

## Addresses

### `getAddresses(input)`

Get addresses from `input`.

---

### `streetAddress(address1,address2,city,state,country,zip)`

Create a location value from a specific street address.

## Maps

### `openInMaps(location)`

Open `location` in the maps app.

## Routing

### `getHalfwayPoint(firstLocation,secondLocation)`

Get the halfway point between `firstLocation` and `secondLocation`.

## Weather

### `getCurrentWeather()`

Get current weather conditions for the users current location.

---

### `getCurrentWeatherAt(location)`

Get current weather conditions at `location`.

---

### `getWeatherDetail(weather,detail)`

Get `detail` of `weather`.

---

### `getWeatherForcast(type)`

Get `type` of forecast for the users current location.

#### Type

- Daily
- Hourly

---

### `getWeatherForcastAt(type,location)`

Get `type` of forecast for `location`.
