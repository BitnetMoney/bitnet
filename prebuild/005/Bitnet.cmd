::fBE1pAF6MU+EWHreyHcjLQlHcBSDKWe1SLcd+uf17O+Lq3E5W+EqdYrVzqeyB+kHwlDmSZQkwnRfjPccHh5aaxunajMi+SAS+DCHC9CVvTOvSEGd6E4kew==
::fBE1pAF6MU+EWHreyHcjLQlHcBSDKWe1SLcd+uf17O+Lq3E5W+EqdYrVzqeyB+kHwlDmSZQkwnRfjPccHh5aaxunajMi+SAS+DCHC9CVvTOvSUeH4EI3ew==
::fBE1pAF6MU+EWHreyHcjLQlHcBSDKWe1SLcd+uf17O+Lq3E5W+EqdYrVzqeyB+kHwlDmSZQkwnRfjPccHh5aaxunajMi+SAS+DCHC9CVvTPjQ1yH419+Hn1x5w==
::YAwzoRdxOk+EWAjk
::fBw5plQjdCyDJGyX8VAjFAtVWQiNcmm7FLoS6+335tajrU4IWecxbJzn/b2aCPUR1kftYZgowjRTm8Rs
::YAwzuBVtJxjWCl3EqQJgSA==
::ZR4luwNxJguZRRnk
::Yhs/ulQjdF+5
::cxAkpRVqdFKZSjk=
::cBs/ulQjdF+5
::ZR41oxFsdFKZSDk=
::eBoioBt6dFKZSDk=
::cRo6pxp7LAbNWATEpSI=
::egkzugNsPRvcWATEpSI=
::dAsiuh18IRvcCxnZtBJQ
::cRYluBh/LU+EWAnk
::YxY4rhs+aU+IeA==
::cxY6rQJ7JhzQF1fEqQJgZksaGAbi
::ZQ05rAF9IBncCkqN+0xwdVsEAlTMaGna
::ZQ05rAF9IAHYFVzEqQISIQ9aSRDi
::eg0/rx1wNQPfEVWB+kM9LVsJDGQ=
::fBEirQZwNQPfEVWB+kM9LVsJDGQ=
::cRolqwZ3JBvQF1fEqQJQ
::dhA7uBVwLU+EWDk=
::YQ03rBFzNR3SWATElA==
::dhAmsQZ3MwfNWATElA==
::ZQ0/vhVqMQ3MEVWAtB9wSA==
::Zg8zqx1/OA3MEVWAtB9wSA==
::dhA7pRFwIByZRRnk
::Zh4grVQjdCyDJGyX8VAjFAtVWQiNcmm7FLoS6+335tajrU4IWecxbJzn/b2aCPUR1kftYdgozn86
::YB416Ek+ZG8=
::
::
::978f952a14a936cc963da21a135fa983
@echo off

:: Copyright 2023 Bitnet
:: This file is part of the Bitnet library.
::
:: This software is provided "as is", without warranty of any kind,
:: express or implied, including but not limited to the warranties
:: of merchantability, fitness for a particular purpose and
:: noninfringement. In no even shall the authors or copyright
:: holders be liable for any claim, damages, or other liability,
:: whether in an action of contract, tort or otherwise, arising
:: from, out of or in connection with the software or the use or
:: other dealings in the software.

title Bitnet Node
color 07

:: Register and initiates the mainnet genesis file.
echo Initiating mainnet genesis...
.\bitnet --datadir bitnet.db init .genesis

:: Starts the node using the parameters specified.
echo Starting node...
.\bitnet --networkid 210 --config .nodeconfig