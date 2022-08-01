import 'package:flutter/material.dart';

import 'colors.dart';

class SchejConstants {
  // [pagePadding] is the amount of padding each page should have from the edge of the screen
  static const EdgeInsets pagePadding =
      EdgeInsets.only(top: 16, left: 16, right: 16);

  static final BorderRadius borderRadius = BorderRadius.circular(10);

  static final BoxShadow boxShadow = BoxShadow(
    color: SchejColors.pureBlack.withOpacity(0.1),
    // color: SchejColors.pureBlack.withOpacity(.9),
    spreadRadius: 0,
    blurRadius: 4,
    offset: const Offset(0, 1),
  );

  static final BoxDecoration listTileDecoration = BoxDecoration(
    color: SchejColors.white,
    borderRadius: SchejConstants.borderRadius,
    border: Border.all(
      color: SchejColors.offWhite,
      width: 1,
    ),
    boxShadow: [SchejConstants.boxShadow],
  );
}
