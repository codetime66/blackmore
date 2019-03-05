#!/bin/sh
curl -d '{"api":"v1","linkseller":{"person":{"type":"Doe (%s)","document":"Jane (%s)"},"machine":{"modelcode":"432","seriesnumber":"Jane (%s)"},"order":{"ordercode":"92383"}}}' -H "Content-Type: application/json" -X POST http://localhost:8080/v1/linkseller


