import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

import '../constants/colors.dart';
import '../constants/constants.dart';
import '../constants/fonts.dart';

class FriendRequestCard extends StatelessWidget {
  final String name;
  final DateTime requestTimestamp;

  const FriendRequestCard({
    Key? key,
    required this.name,
    required this.requestTimestamp,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: SchejColors.white,
        borderRadius: SchejConstants.borderRadius,
        border: Border.all(
          color: SchejColors.offWhite,
          width: 1,
        ),
        boxShadow: [
          BoxShadow(
            color: SchejColors.pureBlack.withOpacity(0.1),
            spreadRadius: 0,
            blurRadius: 4,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      child: ListTile(
        dense: true,
        leading: const CircleAvatar(
          backgroundImage: NetworkImage(
            'https://pbs.twimg.com/media/D8dDZukXUAAXLdY.jpg',
          ),
        ),
        title: Text(
          name,
          style: SchejFonts.subtitle.copyWith(color: SchejColors.pureBlack),
        ),
        subtitle: _buildRequestTimestamp(),
        trailing: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            IconButton(
              icon: const Icon(MdiIcons.close),
              color: SchejColors.black,
              onPressed: () {},
            ),
            IconButton(
              icon: const Icon(MdiIcons.check),
              color: SchejColors.darkGreen,
              onPressed: () {},
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRequestTimestamp() {
    final dateString = DateFormat.MMMMd().format(requestTimestamp);
    final timestampString = 'Requested on $dateString';

    return Text(timestampString,
        style: SchejFonts.body.copyWith(color: SchejColors.pureBlack));
  }
}
