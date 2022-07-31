import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:intl/intl.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

class FriendRequestCard extends StatelessWidget {
  final String name;
  final String pic;
  final DateTime requestTimestamp;

  const FriendRequestCard({
    Key? key,
    required this.name,
    this.pic = 'https://pbs.twimg.com/media/D8dDZukXUAAXLdY.jpg',
    required this.requestTimestamp,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: SchejConstants.listTileDecoration,
      child: ListTile(
        dense: true,
        leading: CircleAvatar(
          backgroundImage: NetworkImage(pic),
        ),
        title: Text(name, style: SchejFonts.subtitle),
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
