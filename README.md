# H3VR.BT2TS
This project will restructure any BoneTome mod to the best of it's ability to facilitate migrating to Thunderstore.

## Why use it?
- You don't need to be familiar with any specific structure.
- User-friendly structures.
- Automatically link to the latest relevant dependency on Thunderstore.
    - Deli, Sideloader, and OtherLoader; but only when required.
- Pulls as much information as possible directly from BT, reducing the amount of information you need to enter yourself.
- Generates (or at least tries to generate) a README.md file based on a description and changelog fields on BT.
- Converts mod names to Thunderstore-supported equivalents.

## How do I use it?
### To download
Either:
- Download the latest EXE directly from: https://github.com/ebkr/H3VR.BT2TS/releases
- Or, build it from source. The executables were built using Go 1.15.7.

### To run
1. Start the executable.
2. Ensure you're connected to the internet.
3. Copy and paste the URL from a BoneTome mod page.
4. Enter a description (up to 250 characters).
5. Enter a website URL (Can be blank, a GitHub URL, the BT mod page, or something else).
6. Wait for the files to generate.
7. Select an option to either leave as-is, or zip for mod manager importing.

## How do I upload to Thunderstore?
The restructured folder will be generated in the same directory as the executable.

You will need to supply an `icon.png` with a resolution of 256px by 256px.

1. Place the `icon.png` in the restructured folder, alongside the `manifest.json` and `README.md`.
2. Zip the files so that `icon.png`, `manifest.json`, and `README.md` are in the root of the zip file.
3. Login to Thunderstore.
4. Upload the zip.

